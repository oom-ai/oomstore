fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .build_server(false)
        .build_client(true)
        .compile_well_known_types(true)
        .compile(&["oomagent.proto", "status.proto"], &["proto"])?;
    Ok(())
}
